# generate_stress_data.py
import os
import time

def generate_dummy_nodes(target_dir="content/hs-code/dummy", num_files=100000):
    """
    Hugo 엔진의 렌더링 한계를 타격하기 위한 물리량 더미 생성기.
    """
    os.makedirs(target_dir, exist_ok=True)
    start_time = time.time()

    # Hugo Frontmatter 파싱 부하를 모사하기 위한 최소/필수 템플릿
    template = """---
title: "HS Code {code}"
date: 2026-04-01
taxRate: {tax}
vatRate: 10
layout: "single"
---
HS Code {code} Landed Cost Calculator Instance.
"""
    
    print(f"[{time.time() - start_time:.2f}s] ⚡ 물리적 스트레스 테스트 시작: {num_files}개 노드 생성 중...")

    for i in range(num_files):
        # 000000 부터 099999 까지 6자리 HS 코드 모사
        code = f"{i:06d}"
        tax = (i % 30) + 5 # 5% ~ 34% 가상 관세율
        
        file_path = os.path.join(target_dir, f"{code}.md")
        
        # 순수 I/O 속도 극대화를 위한 버퍼 라이팅
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(template.format(code=code, tax=tax))

        if i % 25000 == 0 and i > 0:
            print(f"[{time.time() - start_time:.2f}s] 진행률: {i}개 노드 생성 완료...")

    print(f"[{time.time() - start_time:.2f}s] ✅ 테스트 데이터 준비 완료. 총 {num_files}개의 파일이 {target_dir}에 적재됨.")

if __name__ == "__main__":
    # 초기 한계 타격 목표: 100,000 페이지
    generate_dummy_nodes(num_files=100000)
