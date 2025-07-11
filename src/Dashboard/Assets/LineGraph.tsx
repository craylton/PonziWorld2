import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface ChartDataPoint {
  day: number;
  value: number;
}

interface LineGraphProps {
  data: ChartDataPoint[];
  title: string;
  formatTooltip: (value: number) => string;
  formatYAxisTick: (value: number) => string;
  height?: string;
  color?: string;
}

export default function LineGraph({
  data,
  title,
  formatTooltip,
  formatYAxisTick,
  height = '300px',
  color = '#2563eb'
}: LineGraphProps) {
  const chartData = {
    labels: data.map(d => `Day ${d.day}`),
    datasets: [
      {
        label: title,
        data: data.map(d => d.value),
        borderColor: color,
        backgroundColor: `${color}1A`, // Add transparency (10% opacity)
        borderWidth: 2,
        pointBackgroundColor: color,
        pointBorderColor: color,
        pointRadius: 3,
        pointHoverRadius: 5,
        tension: 0.1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: (context: { parsed: { y: number } }) => {
            return formatTooltip(context.parsed.y);
          },
        },
      },
    },
    scales: {
      x: {
        grid: {
          color: 'rgba(255, 255, 255, 0.2)',
        },
        ticks: {
          maxTicksLimit: 7,
          color: 'white',
        },
      },
      y: {
        grid: {
          color: 'rgba(255, 255, 255, 0.2)',
        },
        ticks: {
          callback: (value: number | string) => {
            return formatYAxisTick(value as number);
          },
          maxTicksLimit: 6,
          color: 'white',
        },
      },
    },
  };

  return (
    <div style={{ height }}>
      <Line data={chartData} options={options} />
    </div>
  );
}
