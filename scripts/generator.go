package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Task 정의: 워커에게 전달될 원자적 데이터 단위
type Task struct {
	ID         int
	Keyword    string
	TariffRate float64
	VatRate    float64
}

// 워커 수 제한 (ulimit 방어)
const numWorkers = 100
const totalTasks = 100000
const outputDir = "content/calculator"

func main() {
	start := time.Now()
	fmt.Println("🚀 시스템적 천재 아키텍트 타설 엔진 가동 시작...")

	// 1. 디렉토리 사전 타설 (Race Condition 방어)
	preAllocateDirectories()

	// 2. 워커 풀 파이프라인 구축
	jobs := make(chan Task, totalTasks)
	var wg sync.WaitGroup

	// 워커 스폰 (Spawn)
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// 3. 10만 개 작업 큐잉 (메모리에서 채널로 즉시 주입)
	for i := 1; i <= totalTasks; i++ {
		jobs <- Task{
			ID:         i,
			Keyword:    fmt.Sprintf("item-%d", i),
			TariffRate: 8.0,  // 예시 관세율
			VatRate:    10.0, // 예시 부가세율
		}
	}
	close(jobs) // 큐 주입 완료 선언

	// 4. 모든 I/O 작업 완료 대기
	wg.Wait()

	fmt.Printf("✅ 타설 완료: %d개 파일 생성 (소요시간: %v)\n", totalTasks, time.Since(start))
}

// 워커 스레드: 채널에서 작업을 꺼내어 디스크에 기록
func worker(id int, jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// 파일 분산 배치를 위한 해시 계산 (프랙탈 디렉토리)
		hash := getMD5Hash(j.Keyword)
		folder := hash[:2] // 앞 2자리 추출 (256개 폴더로 분산)
		fileName := fmt.Sprintf("%s.md", j.Keyword)
		fullPath := filepath.Join(outputDir, folder, fileName)

		// 최소화된 Frontmatter 강제 주입 (단일 파일 크기 극소화)
		content := fmt.Sprintf(`---
title: "%s"
tariff_rate: %.1f
vat_rate: %.1f
---`, j.Keyword, j.TariffRate, j.VatRate)

		// 디스크 I/O 타격
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			log.Printf("워커 %d - 파일 쓰기 에러: %v\n", id, err)
		}
	}
}

// 256개 (00~ff)의 헥사 디렉토리를 메인 스레드에서 사전 생성
func preAllocateDirectories() {
	_ = os.MkdirAll(outputDir, 0755)
	for i := 0; i < 256; i++ {
		folder := fmt.Sprintf("%02x", i)
		_ = os.MkdirAll(filepath.Join(outputDir, folder), 0755)
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
