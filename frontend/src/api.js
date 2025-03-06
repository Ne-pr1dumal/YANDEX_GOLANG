const API_URL = 'http://localhost:8080/api/v1';

export const calculateExpression = async (expression) => {
  const response = await fetch(`${API_URL}/calculate`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ expression }),
  });
  return response.json();
};

export const getExpressions = async () => {
  const response = await fetch(`${API_URL}/expressions`);
  const data = await response.json();
  return data.expressions;
};
