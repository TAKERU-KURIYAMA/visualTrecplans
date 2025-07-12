import React, { useState } from 'react';

interface BodyPart {
  id: string;
  name: string;
  muscleGroup: string;
  color: string;
  path: string;
}

interface BodyDiagramProps {
  selectedMuscleGroup?: string;
  onMuscleGroupSelect?: (muscleGroup: string) => void;
  workoutData?: Record<string, number>;
  className?: string;
}

export const BodyDiagram: React.FC<BodyDiagramProps> = ({
  selectedMuscleGroup,
  onMuscleGroupSelect,
  workoutData = {},
  className = ""
}) => {
  const [hoveredPart, setHoveredPart] = useState<string | null>(null);

  const bodyParts: BodyPart[] = [
    {
      id: 'chest',
      name: '胸',
      muscleGroup: 'chest',
      color: '#ff6b6b',
      path: 'M120,80 Q160,60 200,80 Q200,120 160,140 Q120,120 120,80'
    },
    {
      id: 'shoulders-left',
      name: '肩（左）',
      muscleGroup: 'shoulders',
      color: '#dda0dd',
      path: 'M80,70 Q100,50 120,70 Q120,90 100,100 Q80,90 80,70'
    },
    {
      id: 'shoulders-right',
      name: '肩（右）',
      muscleGroup: 'shoulders',
      color: '#dda0dd',
      path: 'M200,70 Q220,50 240,70 Q240,90 220,100 Q200,90 200,70'
    },
    {
      id: 'arms-left',
      name: '腕（左）',
      muscleGroup: 'arms',
      color: '#96ceb4',
      path: 'M60,90 Q80,85 85,110 Q85,140 80,160 Q60,155 55,130 Q55,110 60,90'
    },
    {
      id: 'arms-right',
      name: '腕（右）',
      muscleGroup: 'arms',
      color: '#96ceb4',
      path: 'M235,110 Q240,85 260,90 Q265,110 265,130 Q260,155 240,160 Q235,140 235,110'
    },
    {
      id: 'core',
      name: '腹',
      muscleGroup: 'core',
      color: '#ffd93d',
      path: 'M130,140 Q190,140 190,200 Q160,220 130,200 Q130,170 130,140'
    },
    {
      id: 'back',
      name: '背中',
      muscleGroup: 'back',
      color: '#4ecdc4',
      path: 'M120,80 Q160,60 200,80 Q200,120 160,140 Q120,120 120,80'
    },
    {
      id: 'legs-left',
      name: '脚（左）',
      muscleGroup: 'legs',
      color: '#45b7d1',
      path: 'M130,220 Q140,215 145,240 Q145,280 140,320 Q130,315 125,280 Q125,240 130,220'
    },
    {
      id: 'legs-right',
      name: '脚（右）',
      muscleGroup: 'legs',
      color: '#45b7d1',
      path: 'M175,220 Q185,215 195,220 Q195,240 195,280 Q190,315 180,320 Q175,280 175,240 Q175,220 175,220'
    },
    {
      id: 'glutes',
      name: '臀部',
      muscleGroup: 'glutes',
      color: '#ff9999',
      path: 'M130,200 Q190,200 190,230 Q160,235 130,230 Q130,215 130,200'
    }
  ];

  const getPartOpacity = (part: BodyPart) => {
    if (selectedMuscleGroup && part.muscleGroup === selectedMuscleGroup) {
      return 1;
    }
    if (selectedMuscleGroup && part.muscleGroup !== selectedMuscleGroup) {
      return 0.3;
    }
    if (hoveredPart === part.id) {
      return 0.8;
    }
    const workoutCount = workoutData[part.muscleGroup] || 0;
    if (workoutCount > 0) {
      return 0.6 + (workoutCount / 10) * 0.4; // Scale intensity based on workout count
    }
    return 0.4;
  };

  const handlePartClick = (part: BodyPart) => {
    if (onMuscleGroupSelect) {
      onMuscleGroupSelect(part.muscleGroup);
    }
  };

  return (
    <div className={`relative ${className}`}>
      <svg
        viewBox="0 0 320 340"
        className="w-full h-full"
        style={{ maxWidth: '400px', maxHeight: '500px' }}
      >
        {/* Body outline */}
        <path
          d="M160,20 Q200,30 220,60 Q240,90 240,130 Q240,170 230,200 Q220,230 200,250 Q190,280 185,320 Q175,330 165,330 Q155,330 145,320 Q140,280 130,250 Q110,230 100,200 Q90,170 90,130 Q90,90 110,60 Q140,30 160,20"
          fill="none"
          stroke="#e5e7eb"
          strokeWidth="2"
          className="opacity-30"
        />

        {/* Muscle groups */}
        {bodyParts.map((part) => (
          <g key={part.id}>
            <path
              d={part.path}
              fill={part.color}
              opacity={getPartOpacity(part)}
              stroke="#ffffff"
              strokeWidth="1"
              className="cursor-pointer transition-all duration-200 hover:stroke-2"
              onClick={() => handlePartClick(part)}
              onMouseEnter={() => setHoveredPart(part.id)}
              onMouseLeave={() => setHoveredPart(null)}
            />
          </g>
        ))}

        {/* Labels */}
        {hoveredPart && (
          <g>
            {bodyParts
              .filter(part => part.id === hoveredPart)
              .map((part) => (
                <g key={`label-${part.id}`}>
                  <rect
                    x="10"
                    y="10"
                    width={part.name.length * 8 + 20}
                    height="30"
                    fill="rgba(0,0,0,0.8)"
                    rx="4"
                  />
                  <text
                    x="20"
                    y="30"
                    fill="white"
                    fontSize="14"
                    fontWeight="bold"
                  >
                    {part.name}
                  </text>
                </g>
              ))}
          </g>
        )}
      </svg>

      {/* Legend */}
      <div className="mt-4 text-sm text-gray-600">
        <p className="mb-2 font-medium">筋肉部位</p>
        <div className="grid grid-cols-2 gap-2">
          {Array.from(new Set(bodyParts.map(part => part.muscleGroup))).map(muscleGroup => {
            const part = bodyParts.find(p => p.muscleGroup === muscleGroup);
            const workoutCount = workoutData[muscleGroup] || 0;
            return (
              <div
                key={muscleGroup}
                className={`flex items-center cursor-pointer p-2 rounded hover:bg-gray-50 ${
                  selectedMuscleGroup === muscleGroup ? 'bg-blue-50 border border-blue-200' : ''
                }`}
                onClick={() => onMuscleGroupSelect?.(muscleGroup)}
              >
                <div
                  className="w-3 h-3 rounded-full mr-2"
                  style={{ backgroundColor: part?.color }}
                />
                <span className="flex-1">{part?.name}</span>
                {workoutCount > 0 && (
                  <span className="text-xs bg-gray-200 px-2 py-1 rounded">
                    {workoutCount}
                  </span>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};