import matplotlib.pyplot as plt

s = 0.5
uns = s
mass_i = []
mass_s = []
for i in range (1, 1001):
    mass_i.append(i)
    uns /= 2**i 
    s += uns
    mass_s.append(s)
plt.figure(figsize=(10, 6))
plt.plot(mass_i, mass_s, color='blue', linewidth=2)
plt.title('Зависимость s от i', fontsize=14)
plt.xlabel('i', fontsize=12)
plt.ylabel('s', fontsize=12)
plt.grid(True, linestyle='--', alpha=0.7)
plt.xticks(mass_i)
plt.ylim(0.4, 0.9)
plt.show()