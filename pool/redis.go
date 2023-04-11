package pool

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

var (
	ZSetKey = "_tool_ippool_ZSet_Key"
)

// AddDataToZSet 新增或更新值
func AddDataToZSet(key, member string, score float64) error {
	data := redis.Z{
		Score:  score,
		Member: member,
	}
	_, err := RedisClient.ZAdd(context.Background(), key+ZSetKey, &data).Result()
	if err != nil {
		return err
	}
	if err == redis.Nil {
		return errors.New("AddDataToZSet err:not find:" + key)
	}
	return nil
}

// IncrData 增加使用次数
func IncrData(key, member string, incr float64) error {
	_, err := RedisClient.ZIncrBy(context.Background(), key, incr, member).Result()
	if err != nil {
		return err
	}
	if err == redis.Nil {
		return errors.New("IncrData err:not find:" + key)
	}
	return nil
}

// GetTop1Ip 获取最小使用次数的ip,并将使用次数+1
func GetTop1Ip(key string) (string, error) {
	key = key + ZSetKey
	res, err := RedisClient.ZRange(context.Background(), key, 0, int64(QuarterIpNumber)).Result()

	if err != nil {
		return "", err
	} else if err == redis.Nil {
		return "", errors.New("GetTop1Ip err:not find:" + key)
	} else if res == nil {
		return "", errors.New("redis key is nil")
	}
	mem := res[Random.Intn(QuarterIpNumber)]

	if err := IncrData(key, mem, 1.0); err != nil {
		return "", errors.New("IncrData error" + err.Error())
	}
	return mem, nil
}

// GetAllIp 获取全部Ip
func GetAllIp(key string) (*[]string, error) {
	key = key + ZSetKey
	//获取全部Ip
	res, err := RedisClient.ZRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	} else if err == redis.Nil {
		return nil, errors.New("GetTop1Ip err:not find:" + key)
	} else if res == nil {
		return nil, errors.New("redis key is nil")
	}
	resStr := make([]string, len(res))

	for i, e := range res {
		resStr[i] = e
	}
	return &resStr, nil

}
