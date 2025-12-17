import React, { useState } from 'react';
import { Button, Alert } from '@mui/material';
import { exportDefectsCSV } from '../api/Projects';

const ExportDefectsButton = ({ projectId }) => {
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleExport = async () => {
    try {
      const response = await exportDefectsCSV(projectId);
      setSuccess(response.message);
      setError('');
      setTimeout(() => setSuccess(''), 3000);
    } catch (err) {
      setError(err.message || 'Failed to export CSV');
      setSuccess('');
    }
  };

  return (
    <div>
      <Button
        variant="contained"
        color="primary"
        onClick={handleExport}
      >
        Экспортировать
      </Button>
      {error && (
        <Alert severity="error" sx={{ mt: 2 }}>
          {error}
        </Alert>
      )}
      {success && (
        <Alert severity="success" sx={{ mt: 2 }}>
          {success}
        </Alert>
      )}
    </div>
  );
};

export default ExportDefectsButton;