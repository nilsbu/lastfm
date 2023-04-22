import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import Table from './Table';
import Menu from './Menu';
import './Page.css';
import { buttons, getMenus } from './menus';

// type that we get as JSON. There is more because it's also used for the chart.
interface JSONData {
  data: {
    labels: string[],
    datasets: {
      data: number[]
    }[]
  }
}

function Page() {
  const [method, setMethod] = useState([buttons['topLevel'].buttons[0].function]);

  const getMethod = (methodArray : string[]) => {
    return methodArray.join('/');
  };

  const handleMethodChange = (newMethod : string, index : number) => {
    console.log(`Changing method to ${newMethod} at index ${index}`);
    if (newMethod !== method[index]) {
      var newMethodArray = [...method]; // create a copy of the method array
      if (index === 0) {
        // if the top level method has changed, reset the rest of the method array
        newMethodArray = getMenus(newMethod).map(menu => buttons[menu].default);
      }
      newMethodArray[index] = newMethod;
      console.log(`New method array: ${newMethodArray}`);
      setMethod(newMethodArray); // update the method state with the new array
      fetchData(newMethodArray); // fetch new data
    }
  };

  const transformData = (data : JSONData) => {
    return data.data.labels.map((label, index) => {
      const value = data.data.datasets[0].data[index];
      return { label, value };
    });
  };

  const isFirstRender = useRef(true); // add a ref to keep track of initial render

  useEffect(() => {
    if (isFirstRender.current) { // check if it's the first render
      isFirstRender.current = false;
    } else {
      fetchData(method);
    }
  }, []); // no dependencies, so it only runs once

  const [data, setData] = useState<TableData>([]);

  const fetchData = (method : string[]) => {
    const name = getMethod(method);
    console.log(`Fetching data from http://${window.location.hostname}:3001/json/print/${name}`);
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
          {getMenus(method[0]).map((menu, index) => (
            <Menu
              key={menu}
              onMethodChange={newMethod => handleMethodChange(newMethod, index)}
              buttons={buttons[menu]}
            />
          ))}
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
