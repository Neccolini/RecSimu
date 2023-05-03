import json
import random
import sys

def select_random_numbers(n: int, m: int):
    return random.sample(range(1, n), m)

file_path = sys.argv[1]
steps = int(sys.argv[2])

with open(file_path, 'r') as f:
    data = json.load(f)

adj = data["adjacencies"]
end = len(adj)

data["reconfigure"] = []
cur_cycle = 10000
for i in range(5):
  all_nodes = [_num for _num in range(1, end)]
  rem_nodes = select_random_numbers(end, steps)
  for node in all_nodes:
    data["reconfigure"].append(
      {
        "id": str(node),
        "cycle": cur_cycle,
        "operation": "remove"
      }
    )
    
    if node in rem_nodes:
      data["reconfigure"].append(
        {
          "id": str(node),
          "cycle": cur_cycle + 3000,
          "operation": "rejoin",
          "node_type": "Router",
          "adjacencies": adj[str(node)],
        }
      )
    else:
      data["reconfigure"].append(
        {
          "id": str(node),
          "cycle": cur_cycle + 10,
          "operation": "rejoin",
          "node_type": "Router",
          "adjacencies": adj[str(node)],
        }
      )
  cur_cycle += 4000

with open(file_path, "w") as outfile:
    json.dump(data, outfile)