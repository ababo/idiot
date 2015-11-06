package main

import (
	"hash/fnv"
	"strconv"
	"sync"
)

type cacheBank struct {
	sync.RWMutex
	data map[uint64][]ParseMatch
}

var cacheBanks []cacheBank

func InitCache(bankCount uint) {
	banks := make([]cacheBank, bankCount)
	for i := range banks {
		banks[i].data = make(map[uint64][]ParseMatch)
	}
	cacheBanks = banks
}

func FinalizeCache() {
	cacheBanks = nil
}

func getCacheHash(text, nonterminal string, hypothesesLimit uint) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(text))
	hash.Write([]byte(nonterminal))
	hash.Write([]byte(strconv.Itoa(int(hypothesesLimit))))
	return hash.Sum64()
}

func FindInCache(text, nonterminal string,
	hypothesesLimit uint) ([]ParseMatch, bool) {

	hash := getCacheHash(text, nonterminal, hypothesesLimit)
	bank := &cacheBanks[hash%uint64(len(cacheBanks))]
	bank.RLock()
	matches, ok := bank.data[hash]
	bank.RUnlock()

	return matches, ok
}

func AddToCache(text, nonterminal string,
	hypothesesLimit uint, matches []ParseMatch) {

	hash := getCacheHash(text, nonterminal, hypothesesLimit)
	bank := &cacheBanks[hash%uint64(len(cacheBanks))]
	bank.Lock()
	bank.data[hash] = matches
	bank.Unlock()
}

func ClearCache() {
	InitCache(uint(len(cacheBanks)))
}
