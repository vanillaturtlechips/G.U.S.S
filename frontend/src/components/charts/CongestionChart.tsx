import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend,
  ChartData,
  ChartOptions
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend
);

interface CongestionChartProps {
  gymId: number;
  color?: string;
}

const CongestionChart: React.FC<CongestionChartProps> = ({ gymId, color = '#4ade80' }) => {
  const [chartData, setChartData] = useState<ChartData<'line'>>({
    labels: [],
    datasets: [],
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`/api/gyms/${gymId}/congestion`);
        const result = await response.json();

        // [수정] 초(second) 단위를 추가하여 중복 라벨 방지 (차트 렌더링 정상화)
        const now = new Date().toLocaleTimeString('ko-KR', { 
          hour: '2-digit', 
          minute: '2-digit',
          second: '2-digit' 
        });

        setChartData((prev) => {
          const newLabels = [...(prev.labels || []), now].slice(-10);
          const oldData = prev.datasets[0]?.data || [];
          const newData = [...oldData, result.ratio * 100].slice(-10);

          return {
            labels: newLabels,
            datasets: [
              {
                fill: true,
                label: '혼잡도 (%)',
                data: newData as number[],
                borderColor: color,
                backgroundColor: `${color}33`,
                tension: 0.4,
                pointRadius: 4,
                pointBackgroundColor: color,
              },
            ],
          };
        });
      } catch (err) {
        console.error("데이터 로드 실패:", err);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, [gymId, color]);

  const options: ChartOptions<'line'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (context) => `혼잡도: ${context.parsed.y}%`
        }
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
        grid: { color: 'rgba(128, 128, 128, 0.1)' },
        ticks: { font: { size: 10 } },
      },
      x: {
        grid: { display: false },
        ticks: { font: { size: 10 } },
      },
    },
  };

  return (
    <div className="congestion-chart-container" style={{ width: '100%', height: '100%', minHeight: '200px' }}>
      <Line data={chartData} options={options} />
    </div>
  );
};

export default CongestionChart;