import React from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import Table from './Table';
import Menu from './Menu';
import './Page.css';

function Page() {
    const chartData1 = Array.from({ length: 150 }, (_, i) => ({
        label: `Label ${i + 1}`,
        value: (i + 1) * 10,
      }));
      

    const chartData2 = Array.from({ length: 200 }, (_, i) => ({
    label: `Label ${String.fromCharCode(65 + i)}`,
    value: (i + 1) * 50,
    }));
    

  return (
    <Container fluid>
      <Row>
        <Col>
          <Menu />
        </Col>
      </Row>
      <Row>
        <Col className="table-container">
          <h2>Chart 1</h2>
          <Table data={chartData1} />
        </Col>
        <Col className="table-container">
          <h2>Chart 2</h2>
          <Table data={chartData2} />
        </Col>
      </Row>
    </Container>
  );
}

export default Page;
