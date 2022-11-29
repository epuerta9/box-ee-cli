import os
import random

for i in range(40):
    num = random.randrange(10000,11000)
    os.system(f"boxee tracking add --tracking-number {num}")