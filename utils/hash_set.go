// Copyright 2017 guangbo. All rights reserved.
package utils

import (
	"log"
	"strconv"
	"strings"
	"sync"
)

type HashSet struct {
	Mutex sync.Mutex
	Data  map[string]map[string]string
}

func NewHashSet() *HashSet {
	return &HashSet{
		Data: make(map[string]map[string]string),
	}
}

// KeyPrefix 找到前缀为prefix的key列表
func (s *HashSet) KeyPrefix(prefix string) []string {
	var ret []string
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	for k, _ := range s.Data {
		if strings.HasPrefix(k, prefix) {
			ret = append(ret, k)
		}
	}
	return ret
}

func (s *HashSet) Del(key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Data, key)

	DelKey(key)
}

func (s *HashSet) Hget(key, field string) string {
	if key == "" || field == "" {
		return ""
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	keyInfo, ok := s.Data[key]
	if !ok {
		return ""
	}

	var value string
	value, ok = keyInfo[field]
	if !ok {
		return ""
	}

	return value
}

func (s *HashSet) Load() {
	kfv := GetAllKfv()
	for i := 0; i < len(kfv); i++ {
		log.Println("load set, ", kfv[i].Key, kfv[i].Field, kfv[i].Value)
		s.hset(kfv[i].Key, kfv[i].Field, kfv[i].Value)
	}
}

func (s *HashSet) Hgetall(key string) []string {
	if key == "" {
		return []string{}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	keyInfo, ok := s.Data[key]
	if !ok {
		return []string{}
	}

	var arr []string
	for field, value := range keyInfo {
		arr = append(arr, field)
		arr = append(arr, value)
	}

	return arr
}

func (s *HashSet) Hset(key, field, value string) bool {
	status, ok := s.hset(key, field, value)
	if !ok {
		return false
	}
	if status == 0 {
		UpdateKvs(key, field, value)
	} else if status == 1 {
		AddKvs(key, field, value)
	} else {
		DelKvs(key, field)
	}
	return true
}

func (s *HashSet) hset(key, field, value string) (int, bool) {
	if key == "" || field == "" {
		return 0, false
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	status := 0 //0-update 1-add 2-del
	keyInfo, ok := s.Data[key]
	if !ok {
		keyInfo = make(map[string]string)
		s.Data[key] = keyInfo
	}

	_, ok = keyInfo[field]

	if value == "" {
		delete(keyInfo, field)
		status = 2
	} else if ok {
		keyInfo[field] = value
		status = 0
	} else {
		keyInfo[field] = value
		status = 1
	}

	return status, true
}

func (s *HashSet) Hincrby(key, field string, add int) string {
	status, value, ok := s.hincrby(key, field, add)
	if !ok {
		return ""
	}
	if status == 0 {
		UpdateKvs(key, field, value)
	} else {
		AddKvs(key, field, value)
	}
	return value
}

func (s *HashSet) hincrby(key, field string, add int) (int, string, bool) {
	if key == "" || field == "" {
		return 0, "", false
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	status := 0 //0-update 1-add 2-del
	keyInfo, ok := s.Data[key]
	if !ok {
		keyInfo = make(map[string]string)
		s.Data[key] = keyInfo
	}

	var value string
	value, ok = keyInfo[field]

	oldValue, _ := strconv.Atoi(value)
	keyInfo[field] = strconv.Itoa(add + oldValue)
	if ok {
		status = 0
	} else {
		status = 1
	}

	return status, keyInfo[field], true
}
