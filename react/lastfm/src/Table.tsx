import React from 'react';
import { Table as BootstrapTable } from 'react-bootstrap';

function Table(props : { data: { label: string, value: number }[] }) {
  const { data } = props;

  if (!data) {
    return null;
  }

  console.log(data);

  return (
    <BootstrapTable striped bordered hover>
      <thead>
        <tr>
          <th>Label</th>
          <th>Value</th>
        </tr>
      </thead>
      <tbody>
        {data.map(item => (
          <tr key={item.label}>
            <td>{item.label}</td>
            <td>{item.value}</td>
          </tr>
        ))}
      </tbody>
    </BootstrapTable>
  );
}

export default Table;
