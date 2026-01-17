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

// ChartJS 플러그인 등록
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
  gymId: number; // 지점 ID
  color?: string; // 기존 디자인 포인트 컬러를 받아올 수 있게 확장
}

const CongestionChart: React.FC<CongestionChartProps> = ({ gymId, color = '#4ade80' }) => {
  // 차트 데이터 상태 관리 (타입 지정)
  const [chartData, setChartData] = useState<ChartData<'line'>>({
    labels: [],
    datasets: [],
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(`/api/gyms/${gymId}/congestion`);
        const result = await response.json();

        const now = new Date().toLocaleTimeString('ko-KR', { 
          hour: '2-digit', 
          minute: '2-digit' 
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
                backgroundColor: `${color}33`, // 컬러에 투명도 20% 추가
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

  // 차트 옵션 설정 (타입 지정)
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