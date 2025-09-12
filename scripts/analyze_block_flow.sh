#!/bin/bash

# Script to analyze the full block processing flow timing
# ProcessProposal -> State Transition -> FinalizeBlock

echo "==================================================================="
echo "       BeaconKit Block Processing Flow Performance Analysis"
echo "==================================================================="
echo ""

if [ -z "$1" ]; then
    echo "Usage: $0 <log_file>"
    echo "Example: $0 beacond.log"
    exit 1
fi

LOG_FILE="$1"

if [ ! -f "$LOG_FILE" ]; then
    echo "Error: Log file $LOG_FILE not found"
    exit 1
fi

echo "Analyzing log file: $LOG_FILE"
echo ""

# Function to calculate statistics
calc_stats() {
    awk '{
        sum+=$1; 
        sumsq+=$1*$1; 
        count++; 
        if(NR==1 || $1<min) min=$1; 
        if(NR==1 || $1>max) max=$1
    } 
    END {
        if(count>0) {
            avg=sum/count;
            if(count>1) {
                variance=(sumsq-sum*sum/count)/(count-1);
                if(variance>0) stddev=sqrt(variance); else stddev=0;
            } else {
                stddev=0;
            }
            printf "  Samples: %d\n", count;
            printf "  Average: %.2f ms\n", avg;
            printf "  Min: %.2f ms\n", min;
            printf "  Max: %.2f ms\n", max;
            printf "  StdDev: %.2f ms\n", stddev;
        } else {
            print "  No data found"
        }
    }'
}

echo "==================================================================="
echo "   Time from ProcessProposal START to State Transition START"
echo "==================================================================="
echo ""
echo "This measures the preparation time before state transition begins."
echo "With optimistic payload building, this includes pre-fetching operations."
echo ""

echo "WITH Optimistic Payload Building (optimistic_enabled=true):"
echo "-----------------------------------------------------------"
grep "\[BENCHMARK\] ProcessProposal -> State Transition START" "$LOG_FILE" | \
    grep "optimistic_enabled=true" | \
    awk -F'time_before_transition_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "WITHOUT Optimistic Payload Building (optimistic_enabled=false):"
echo "----------------------------------------------------------------"
grep "\[BENCHMARK\] ProcessProposal -> State Transition START" "$LOG_FILE" | \
    grep "optimistic_enabled=false" | \
    awk -F'time_before_transition_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "==================================================================="
echo "                State Transition Duration (ProcessProposal)"
echo "==================================================================="
echo ""

echo "WITH Optimistic Payload Building:"
echo "----------------------------------"
grep "\[BENCHMARK\] State Transition COMPLETE (ProcessProposal)" "$LOG_FILE" | \
    grep "optimistic_enabled=true" | \
    awk -F'duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "WITHOUT Optimistic Payload Building:"
echo "-------------------------------------"
grep "\[BENCHMARK\] State Transition COMPLETE (ProcessProposal)" "$LOG_FILE" | \
    grep "optimistic_enabled=false" | \
    awk -F'duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "==================================================================="
echo "            Total ProcessProposal Duration"
echo "==================================================================="
echo ""

echo "WITH Optimistic Payload Building:"
echo "----------------------------------"
grep "\[BENCHMARK\] ProcessProposal FULL COMPLETE" "$LOG_FILE" | \
    grep "optimistic_enabled=true" | \
    awk -F'total_duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "WITHOUT Optimistic Payload Building:"
echo "-------------------------------------"
grep "\[BENCHMARK\] ProcessProposal FULL COMPLETE" "$LOG_FILE" | \
    grep "optimistic_enabled=false" | \
    awk -F'total_duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "==================================================================="
echo "                State Transition Duration (FinalizeBlock)"
echo "==================================================================="
echo ""

echo "WITH Optimistic Payload Building:"
echo "----------------------------------"
grep "\[BENCHMARK\] State Transition COMPLETE (FinalizeBlock)" "$LOG_FILE" | \
    grep "optimistic_enabled=true" | \
    awk -F'duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "WITHOUT Optimistic Payload Building:"
echo "-------------------------------------"
grep "\[BENCHMARK\] State Transition COMPLETE (FinalizeBlock)" "$LOG_FILE" | \
    grep "optimistic_enabled=false" | \
    awk -F'duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "==================================================================="
echo "                    FinalizeBlock Total Duration"
echo "==================================================================="
echo ""

echo "WITH Optimistic Payload Building:"
echo "----------------------------------"
grep "\[BENCHMARK\] FinalizeBlock COMPLETE" "$LOG_FILE" | \
    grep "optimistic_enabled=true" | \
    awk -F'total_duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "WITHOUT Optimistic Payload Building:"
echo "-------------------------------------"
grep "\[BENCHMARK\] FinalizeBlock COMPLETE" "$LOG_FILE" | \
    grep "optimistic_enabled=false" | \
    awk -F'total_duration_ms=' '{print $2}' | awk '{print $1}' | \
    calc_stats

echo ""
echo "==================================================================="
echo "                           KEY INSIGHTS"
echo "==================================================================="
echo ""
echo "What to look for:"
echo ""
echo "1. Time before State Transition (ProcessProposal -> State Transition):"
echo "   - Should be HIGHER with optimistic=true (due to pre-fetch operations)"
echo "   - But this extra time is spent doing useful work in advance"
echo ""
echo "2. State Transition Duration:"
echo "   - Should be LOWER with optimistic=true in ProcessProposal"
echo "   - Because state operations were pre-computed during pre-fetch"
echo ""
echo "3. Total ProcessProposal Duration:"
echo "   - May be similar or slightly higher with optimistic=true"
echo "   - But the work is better distributed (less blocking)"
echo ""
echo "4. Overall Block Commit Speed:"
echo "   - Measured from ProcessProposal start to FinalizeBlock complete"
echo "   - Should show improvement with optimistic payload building"
echo ""
echo "Testing Instructions:"
echo "1. Run node with: --beacon-kit.validator.enable-optimistic-payload-builds=true"
echo "2. Collect logs for ~100 blocks"
echo "3. Run node with: --beacon-kit.validator.enable-optimistic-payload-builds=false"
echo "4. Collect logs for ~100 blocks"
echo "5. Compare the metrics using this script"
echo ""