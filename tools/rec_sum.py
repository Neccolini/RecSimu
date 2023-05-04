import re
import sys

mean_pattern = r"mean reconfiguration latency\s+(\d+(?:\.\d+)?)\s+\[cycle\]"
max_pattern =  r"max reconfiguration latency\s+(\d+(?:\.\d+)?)\s+\[cycle\]"

file_path = sys.argv[1]
with open(file_path, "r") as file:
    data = file.read()

mean_matches = re.findall(mean_pattern, data)
max_matches =  re.findall(max_pattern, data)
mean = sum(list(map(float, mean_matches)))/len(mean_matches)
max_mean = sum(list(map(float, max_matches)))/len(max_matches)
max_max = max(list(map(float, max_matches)))
print("mean: ", mean, "max: ", max_mean, "max of max: ", max_max)
