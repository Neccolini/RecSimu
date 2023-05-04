import re
import sys
def find_nth_occurrence_and_match_after(filename, s, n, r):
    with open(filename, 'r') as f:
        content = f.read()
        idx = -1
        for i in range(n):
            idx = content.find(s, idx + 1)
            if idx == -1:
                return None
        match = re.search(r, content[idx+len(s):])
        return match.group(0) if match else None
filename = sys.argv[1]
r = r"cycle \d*"
s = "joined Network"
n = int(sys.argv[2])
print(find_nth_occurrence_and_match_after(filename, s, n, r))