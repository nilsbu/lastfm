import React from 'react';
import { Table as BootstrapTable } from 'react-bootstrap';

function Table(props : { data: TableData }) {
  const { data } = props;

  if (!data) {
    return null;
  }

  return (
    <BootstrapTable striped bordered hover>
      <tbody>
        {data.map((item, index) => (
          <tr key={item.label}>
            <td>{index + 1}</td>
            <td>{item.label}</td>
            <td>{item.value.toFixed(2)}</td>
          </tr>
        ))}
      </tbody>
    </BootstrapTable>
  );
}

export default Table;
