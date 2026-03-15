package storage

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	t "github.com/uaysk/souin-redis/configurationtypes"
	"github.com/uaysk/souin-redis/pkg/rfc"
	"github.com/uaysk/souin-redis/pkg/storage/types"
)

type redisStorage struct {
	addr       string
	username   string
	password   string
	clientName string
	db         int
	stale      time.Duration
	timeout    time.Duration
}

type respValue struct {
	kind    byte
	text    string
	bulk    []byte
	array   []respValue
	integer int64
}

func newRedisStorage(c t.AbstractConfigurationInterface) (types.Storer, error) {
	cfg := c.GetDefaultCache().GetRedis()
	st := &redisStorage{
		stale:   c.GetDefaultCache().GetStale(),
		timeout: time.Second,
	}

	if raw, ok := cfg.Configuration.(map[string]interface{}); ok {
		if addr := parseAddress(raw["InitAddress"]); addr != "" {
			st.addr = addr
		}
		if timeout := parseDuration(raw["DialTimeout"]); timeout > 0 {
			st.timeout = timeout
		}
		st.username = parseString(raw["Username"])
		st.password = parseString(raw["Password"])
		st.clientName = parseString(raw["ClientName"])
		st.db = parseInt(raw["SelectDB"])
	}

	if st.addr == "" && cfg.URL != "" {
		st.addr = strings.Split(cfg.URL, ",")[0]
	}

	if st.addr == "" {
		return nil, errors.New("no redis address configured")
	}

	return st, nil
}

func (provider *redisStorage) Name() string {
	return "REDIS"
}

func (provider *redisStorage) Uuid() string {
	return fmt.Sprintf("%s-%d", provider.addr, provider.db)
}

func (provider *redisStorage) ListKeys() []string {
	keys := []string{}
	now := time.Now()

	for _, element := range provider.keys("IDX_*") {
		value := provider.Get(element)
		if len(value) == 0 {
			continue
		}

		mapping, err := rfc.DecodeMapping(value)
		if err != nil {
			continue
		}

		for _, item := range mapping.Mapping {
			if item.FreshTime.Before(now) && item.StaleTime.Before(now) {
				continue
			}

			keys = append(keys, item.RealKey)
		}
	}

	return keys
}

func (provider *redisStorage) MapKeys(prefix string) map[string]string {
	keys := map[string]string{}
	for _, key := range provider.keys(prefix + "*") {
		short, _ := strings.CutPrefix(key, prefix)
		keys[short] = string(provider.Get(key))
	}

	return keys
}

func (provider *redisStorage) Get(key string) []byte {
	resp, err := provider.doCommand([][]byte{[]byte("GET"), []byte(key)})
	if err != nil || resp.kind == '$' && resp.bulk == nil {
		return nil
	}

	return resp.bulk
}

func (provider *redisStorage) Set(key string, value []byte, duration time.Duration) error {
	args := [][]byte{[]byte("SET"), []byte(key), value}
	if duration > 0 {
		args = append(args, []byte("PX"), []byte(strconv.FormatInt(duration.Milliseconds(), 10)))
	}

	_, err := provider.doCommand(args)
	return err
}

func (provider *redisStorage) Delete(key string) {
	_, _ = provider.doCommand([][]byte{[]byte("DEL"), []byte(key)})
}

func (provider *redisStorage) DeleteMany(key string) {
	pattern := key
	if !strings.ContainsAny(pattern, "*?[]") {
		pattern += "*"
	}

	keys := provider.keys(pattern)
	if len(keys) == 0 {
		return
	}

	args := make([][]byte, 0, len(keys)+1)
	args = append(args, []byte("DEL"))
	for _, candidate := range keys {
		args = append(args, []byte(candidate))
	}

	_, _ = provider.doCommand(args)
}

func (provider *redisStorage) Init() error {
	_, err := provider.doCommand([][]byte{[]byte("PING")})
	return err
}

func (provider *redisStorage) Reset() error {
	keys := provider.keys("*")
	if len(keys) == 0 {
		return nil
	}

	args := make([][]byte, 0, len(keys)+1)
	args = append(args, []byte("DEL"))
	for _, key := range keys {
		args = append(args, []byte(key))
	}

	_, err := provider.doCommand(args)
	return err
}

func (provider *redisStorage) GetMultiLevel(key string, req *http.Request, validator *types.Revalidator) (fresh *http.Response, stale *http.Response) {
	value := provider.Get("IDX_" + key)
	if len(value) == 0 {
		return
	}

	fresh, stale, _ = rfc.MappingElection(provider, value, req, validator)
	return
}

func (provider *redisStorage) SetMultiLevel(baseKey, variedKey string, value []byte, variedHeaders http.Header, etag string, duration time.Duration, realKey string) error {
	now := time.Now()

	if err := provider.Set(variedKey, value, duration+provider.stale); err != nil {
		return err
	}

	mappingKey := "IDX_" + baseKey
	current := provider.Get(mappingKey)
	updated, err := rfc.MappingUpdater(variedKey, current, now, now.Add(duration), now.Add(duration+provider.stale), variedHeaders, etag, realKey)
	if err != nil {
		return err
	}

	return provider.Set(mappingKey, updated, 0)
}

func (provider *redisStorage) keys(pattern string) []string {
	resp, err := provider.doCommand([][]byte{[]byte("KEYS"), []byte(pattern)})
	if err != nil {
		return nil
	}

	keys := make([]string, 0, len(resp.array))
	for _, value := range resp.array {
		if value.kind == '$' {
			keys = append(keys, string(value.bulk))
		}
	}

	return keys
}

func (provider *redisStorage) doCommand(args [][]byte) (respValue, error) {
	conn, err := net.DialTimeout("tcp", provider.addr, provider.timeout)
	if err != nil {
		return respValue{}, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(provider.timeout))
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	if err = provider.initializeConnection(rw); err != nil {
		return respValue{}, err
	}

	if err = writeRESPCommand(rw, args); err != nil {
		return respValue{}, err
	}
	if err = rw.Flush(); err != nil {
		return respValue{}, err
	}

	resp, err := readRESPValue(rw.Reader)
	if err != nil {
		return respValue{}, err
	}
	if resp.kind == '-' {
		return respValue{}, errors.New(resp.text)
	}

	return resp, nil
}

func (provider *redisStorage) initializeConnection(rw *bufio.ReadWriter) error {
	if provider.password != "" {
		args := [][]byte{[]byte("AUTH")}
		if provider.username != "" {
			args = append(args, []byte(provider.username))
		}
		args = append(args, []byte(provider.password))
		if err := execInitCommand(rw, args); err != nil {
			return err
		}
	}

	if provider.db != 0 {
		if err := execInitCommand(rw, [][]byte{[]byte("SELECT"), []byte(strconv.Itoa(provider.db))}); err != nil {
			return err
		}
	}

	if provider.clientName != "" {
		if err := execInitCommand(rw, [][]byte{[]byte("CLIENT"), []byte("SETNAME"), []byte(provider.clientName)}); err != nil {
			return err
		}
	}

	return nil
}

func execInitCommand(rw *bufio.ReadWriter, args [][]byte) error {
	if err := writeRESPCommand(rw, args); err != nil {
		return err
	}
	if err := rw.Flush(); err != nil {
		return err
	}

	resp, err := readRESPValue(rw.Reader)
	if err != nil {
		return err
	}
	if resp.kind == '-' {
		return errors.New(resp.text)
	}

	return nil
}

func writeRESPCommand(rw *bufio.ReadWriter, args [][]byte) error {
	if _, err := fmt.Fprintf(rw, "*%d\r\n", len(args)); err != nil {
		return err
	}

	for _, arg := range args {
		if _, err := fmt.Fprintf(rw, "$%d\r\n", len(arg)); err != nil {
			return err
		}
		if _, err := rw.Write(arg); err != nil {
			return err
		}
		if _, err := rw.WriteString("\r\n"); err != nil {
			return err
		}
	}

	return nil
}

func readRESPValue(reader *bufio.Reader) (respValue, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return respValue{}, err
	}

	switch prefix {
	case '+', '-', ':':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}

		resp := respValue{kind: prefix, text: line}
		if prefix == ':' {
			resp.integer, _ = strconv.ParseInt(line, 10, 64)
		}

		return resp, nil
	case '$':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}

		size, err := strconv.Atoi(line)
		if err != nil {
			return respValue{}, err
		}
		if size < 0 {
			return respValue{kind: '$', bulk: nil}, nil
		}

		payload := make([]byte, size+2)
		if _, err = io.ReadFull(reader, payload); err != nil {
			return respValue{}, err
		}

		return respValue{kind: '$', bulk: payload[:size]}, nil
	case '*':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}

		count, err := strconv.Atoi(line)
		if err != nil {
			return respValue{}, err
		}
		if count < 0 {
			return respValue{kind: '*'}, nil
		}

		items := make([]respValue, 0, count)
		for i := 0; i < count; i++ {
			item, err := readRESPValue(reader)
			if err != nil {
				return respValue{}, err
			}
			items = append(items, item)
		}

		return respValue{kind: '*', array: items}, nil
	default:
		buffer := bytes.NewBuffer([]byte{prefix})
		rest, _ := reader.ReadBytes('\n')
		buffer.Write(rest)
		return respValue{}, fmt.Errorf("unsupported redis response: %q", buffer.String())
	}
}

func readRESPLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r"), nil
}

func parseAddress(value interface{}) string {
	switch typed := value.(type) {
	case []interface{}:
		if len(typed) == 0 {
			return ""
		}
		return parseString(typed[0])
	case []string:
		if len(typed) == 0 {
			return ""
		}
		return typed[0]
	default:
		return parseString(value)
	}
}

func parseString(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	case []byte:
		return string(typed)
	default:
		return fmt.Sprint(value)
	}
}

func parseInt(value interface{}) int {
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		i, _ := strconv.Atoi(typed)
		return i
	default:
		return 0
	}
}

func parseDuration(value interface{}) time.Duration {
	switch typed := value.(type) {
	case time.Duration:
		return typed
	case string:
		d, _ := time.ParseDuration(typed)
		return d
	case int:
		return time.Duration(typed)
	case int64:
		return time.Duration(typed)
	case float64:
		return time.Duration(typed)
	default:
		return 0
	}
}
