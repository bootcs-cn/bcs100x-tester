#!/bin/bash
# 批量测试所有 stage 的 solution
# 用法: ./scripts/test-all-solutions.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TESTER_DIR="$(dirname "$SCRIPT_DIR")"
SOLUTION_DIR="${TESTER_DIR}/../bcs100x-solution"

# 构建 tester
cd "$TESTER_DIR"
go build -o bcs100x-tester .

# Stage 列表（按课程顺序）
STAGES=(
    "hello"
    "mario-less"
    "mario-more"
    "cash"
    "credit"
    "scrabble"
    "readability"
    "caesar"
    "substitution"
    "sort"
    "plurality"
    "runoff"
    "tideman"
    "volume"
    "filter-less"
    "filter-more"
    "recover"
    "inheritance"
    "speller"
    "sentimental-hello"
    "sentimental-mario-less"
    "sentimental-mario-more"
    "sentimental-cash"
    "sentimental-credit"
    "sentimental-readability"
    "dna"
    "songs"
    "movies"
    "fiftyville"
    "finance"
)

PASSED=0
FAILED=0
SKIPPED=0
TOTAL_TIME=0

echo "=========================================="
echo "  BCS100X Solution Tester"
echo "=========================================="
echo ""

for stage in "${STAGES[@]}"; do
    stage_dir="${SOLUTION_DIR}/${stage}"
    
    if [ ! -d "$stage_dir" ]; then
        echo "⏭️  [$stage] SKIPPED - directory not found"
        ((SKIPPED++))
        continue
    fi
    
    printf "🧪 %-15s Testing... " "[$stage]"
    
    start_time=$(python3 -c 'import time; print(time.time())')
    
    if ./bcs100x-tester -d="$stage_dir" -s="$stage" > /dev/null 2>&1; then
        end_time=$(python3 -c 'import time; print(time.time())')
        elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
        echo "✅ PASSED (${elapsed}s)"
        ((PASSED++))
    else
        end_time=$(python3 -c 'import time; print(time.time())')
        elapsed=$(python3 -c "print(f'{$end_time - $start_time:.2f}')")
        echo "❌ FAILED (${elapsed}s)"
        ((FAILED++))
    fi
    
    TOTAL_TIME=$(python3 -c "print(f'{$TOTAL_TIME + $elapsed:.2f}')")
done

echo ""
echo "=========================================="
echo "  Results: $PASSED passed, $FAILED failed, $SKIPPED skipped"
echo "  Total time: ${TOTAL_TIME}s"
echo "=========================================="

if [ $FAILED -gt 0 ]; then
    exit 1
fi
