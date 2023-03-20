#!/bin/sh

# 計算問題の入ったファイル
filename="rate"

# 行ごとに読み込む
while read line; do
  # 計算式を評価する
  result=$(echo "scale=10; $line" | bc)
  # 結果を表示する
  echo "$result"
done < "$filename"
