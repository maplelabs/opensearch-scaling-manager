import matplotlib.pyplot as plt

font1 = {'family': 'serif', 'color': 'red', 'size': 15}
font2 = {'family': 'serif', 'color': 'darkred', 'size': 10}

plt.plot(x, y)
# plt.xlim([0, duration_minutes])  # set range for x axis
# plt.ylim([0, 60])  # set range for x axis
plt.xlabel('Time (in minutes) -->', font2)
# plt.ylabel('Ingestion Rate (in GB/hr)', font2)
plt.title('Resultant Load', font1)
plt.grid()

plt.show()

