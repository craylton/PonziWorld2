interface AssetSummaryChartProps {
    historicalValues: number[];
    width?: number;
    height?: number;
}

export default function AssetSummaryChart({ 
    historicalValues, 
    width = 32, 
    height = 20 
}: AssetSummaryChartProps) {
    const min = Math.min(...historicalValues);
    const max = Math.max(...historicalValues);
    const points = historicalValues.map((v, i) => {
        const x = historicalValues.length > 1 ? (i / (historicalValues.length - 1)) * width : 0;
        const y = max === min ? height / 2 : height - ((v - min) / (max - min)) * height;
        return `${x},${y}`;
    }).join(' ');

    return (
        <svg
            width={width}
            height={height}
            viewBox={`0 0 ${width} ${height}`}
            xmlns="http://www.w3.org/2000/svg"
            preserveAspectRatio="none"
            style={{ display: 'block', flex: 'none' }}
        >
            <polyline
                points={points}
                fill="none"
                stroke="#ffffff"
                strokeWidth={2}
                strokeLinecap="round"
            />
        </svg>
    );
}
