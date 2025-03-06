import React, { useState } from 'react';
import { calculateExpression, getExpressions } from './api';
import './index.css';

export default function App() {
  const [expression, setExpression] = useState('');
  const [expressions, setExpressions] = useState([]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const result = await calculateExpression(expression);
    setExpressions(await getExpressions());
    setExpression('');
  };

  return (
    <div className="container">
      <h1>Expression Calculator</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={expression}
          onChange={(e) => setExpression(e.target.value)}
          placeholder="Enter expression (e.g., 2+2*2)"
        />
        <button type="submit">Calculate</button>
      </form>

      <div className="results">
        <h2>Calculations History</h2>
        <table>
          <thead>
            <tr>
              <th>Expression</th>
              <th>Status</th>
              <th>Result</th>
            </tr>
          </thead>
          <tbody>
            {expressions.map((expr) => (
              <tr key={expr.id}>
                <td>{expr.expression}</td>
                <td className={`status ${expr.status}`}>{expr.status}</td>
                <td>{expr.result || '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}