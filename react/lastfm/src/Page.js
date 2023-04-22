import React, { useEffect, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import Table from './Table';
import Menu from './Menu';
import './Page.css';

function Page() {
  // button functionality
  const buttons = [
    { function: 'total', name: 'Total' },
    { function: 'fade/365', name: 'Fade 1y' },
    { function: 'fade/3653', name: 'Fade 10y' },
  ];

  const [method, setMethod] = useState('');

  const handleMethodChange = (newMethod) => {
    if (newMethod !== method) {
      fetchData(newMethod);
    }
    setMethod(newMethod);
  };

  const transformData = (data) => {
    return data.data.labels.map((label, index) => {
      const value = data.data.datasets[0].data[index];
      return { label, value };
    });
  };

  useEffect(() => {
    fetchData(buttons[0].function); // Fetch data once at initialization
  }, []); // empty array as second argument to ensure that useEffect only runs once

  const [data, setData] = useState([]);

  const fetchData = (name) => {
    fetch(`http://${window.location.hostname}:3001/json/print/${name}`)
      .then(response => response.json())
      .then(data => transformData(data))
      .then(data => setData(data))
      .catch(error => console.error(error));
  };

  return (
    <Container fluid>
      <Row>
        <Col>
          <Menu onMethodChange={handleMethodChange} buttons={buttons} />
        </Col>
      </Row>
      <Row>
        <Col className="table-container">
          <Table data={data} />
        </Col>
      </Row>
    </Container>
  );
}

export default Page;
