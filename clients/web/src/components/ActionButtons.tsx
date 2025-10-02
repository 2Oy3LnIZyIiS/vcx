import React from 'react';

interface ActionButtonsProps {
  onHealthCheck: () => void;
  onInit: () => void;
  loading: boolean;
}

export const ActionButtons: React.FC<ActionButtonsProps> = ({ 
  onHealthCheck, 
  onInit, 
  loading 
}) => {
  return (
    <div style={{ margin: '20px' }}>
      <button onClick={onHealthCheck} disabled={loading}>
        Check Agent Health
      </button>
      
      <button onClick={onInit} disabled={loading} style={{ marginLeft: '10px' }}>
        Initialize Project
      </button>
    </div>
  );
};