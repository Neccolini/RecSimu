for i in {1..20}
do
  bash tools/gen.sh
  bash tools/run.sh
  bash tools/summarize_all.sh summarized_$i.txt
done
