#!/bin/bash
echo "🚀 클라우드 타설 터빈 점화 시작..."

# 1. 클라우드 서버의 CPU를 착취하여 10만 개 마크다운 즉석 생산
go run scripts/generator.go

# 2. 0ms 전환 터빈 정적 압축 렌더링 (--minify로 용량 극한 압축)
hugo --minify

echo "✅ 글로벌 엣지 배포 준비 완료."
