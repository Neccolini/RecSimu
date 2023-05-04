import re
import sys

mean_pattern = r"mean reconfiguration latency\s+(\d+(?:\.\d+)?)\s+\[cycle\]"
max_pattern =  r"max reconfiguration latency\s+(\d+(?:\.\d+)?)\s+\[cycle\]"
start_pattern = r"(\d0) (0.0\d)\ntotal packets: (\d+) / (\d+)\n.*\naverage latency: (\d*\.\d*) \[cycle\]"
file_path = sys.argv[1]
with open(file_path, "r") as file:
    data = file.read()


matches = re.findall(start_pattern, data)
hash_map = {}
for tp in matches:
  key = (tp[0], tp[1])
  diff = int(tp[3]) - int(tp[2])
  res = (float(tp[2]) * float(tp[4]) + diff * 4000) / float(tp[3])
  if diff < 400:
    res = float(tp[4])
  
  hash_map[key] = hash_map.get(key, 0.0) + res / 5
print(hash_map)
prev = "10"
for key, val in hash_map.items():
  if prev != key[0]:
    print()
    prev = key[0]
  print(key[1], ", ", val, sep="")