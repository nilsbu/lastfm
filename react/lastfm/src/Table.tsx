import React, { ReactNode } from 'react';
import { Table as BootstrapTable } from 'react-bootstrap';

function Table(props : { data: TableData }) {
  const { data } = props;

  if (!data) {
    return null;
  }

  return (
    <BootstrapTable striped bordered hover>
      <tbody>
        {data.map((item, index) => {
          // Calculate position and value differences if previous values exist
          let posDiff: ReactNode = '';
          let valueDiff: ReactNode = '';
          if(item.prevPos !== undefined) {
            const diff = index + 1 - item.prevPos;
            const diffStr = diff > 0 ? `+${diff}` : diff === 0 ? '=0' : diff;
            const color = diff > 0 ? 'red' : diff === 0 ? 'blue' : 'green';
            posDiff = <span>(<span style={{ color: color }}> {diffStr}</span>)</span>;
          }
          if(item.prevValue !== undefined) {
            const diff = item.value - item.prevValue;
            const diffStr = diff > 0 ? `+${diff.toFixed(2)}` : diff === 0 ? '=0.00' : diff.toFixed(2);
            const tolerance = 0.01;
            const color = diff > tolerance ? 'green' : Math.abs(diff) <= tolerance ? 'blue' : 'red';

            valueDiff = <span>(<span style={{ color: color }}> {diffStr}</span>)</span>;
          }
          
          return (
            <tr key={item.label}>
              <td>{index + 1} {posDiff}</td>
              <td>{item.label}</td>
              <td>{item.value.toFixed(2)} {valueDiff}</td>
            </tr>
          );
        })}
      </tbody>
    </BootstrapTable>
  );
}

export default Table;
