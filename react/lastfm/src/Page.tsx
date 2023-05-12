import React, { useEffect, useRef, useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import Table from './Table';
import Menu from './Menu';
import './Page.css';
import { MenuChoice, menuDefinition, getMenus, getQuery, transformMethod } from './menus';

// type that we get as JSON. There is more because it's also used for the chart.
interface JSONData {
  chart: {
    data: {
      title: string;
      value: number;
      prevPos?: number;
      prevValue?: number;
    }[];
  };
};

function Page() {
  const [method, setMethod] = useState<MenuChoice>({topLevel: 'total', functionParam: '', filter: 'all', filterParam: ''});

  const handleMethodChange = (newMethod : string, index : string) => {
    var newChoice : MenuChoice;
    var param : string;
    if (index === 'topLevel') {
      param = newMethod === 'total' ? '' : menu[newMethod].default;
      newChoice = {topLevel: newMethod, functionParam: param, filter: method.filter, filterParam: method.filterParam};
    } else if (index === 'fade' || index === 'period') {
      newChoice = {topLevel: method.topLevel, functionParam: newMethod, filter: method.filter, filterParam: method.filterParam};
    } else if (index === 'filter') {
      param = newMethod === 'all' ? '' : menu['filter'].default;
      newChoice = {topLevel: method.topLevel, functionParam: method.functionParam, filter: newMethod, filterParam: param};
    } else { // must be filterParam
      newChoice = {topLevel: method.topLevel, functionParam: method.functionParam, filter: method.filter, filterParam: newMethod};
    }

    setMethod(newChoice);
    fetchData(newChoice); // fetch new data
  };

  const transformData = (data : JSONData) => {
    return data.chart.data.map((line) => {
      return { 
        label: line.title,
        value: line.value,
        prevPos: line.prevPos !== undefined ? line.prevPos + 1 : undefined,
        prevValue: line.prevValue !== undefined ? line.prevValue : undefined,
      };
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

  const fetchData = (method : MenuChoice) => {
    const name = getQuery(transformMethod(method));

    const hostName = process.env.NODE_ENV === 'production' ? '' : `http://${window.location.hostname}:3001`;
    console.log(`Fetching data from ${hostName}/json/print/${name}`);
    fetch(`${hostName}/json/print/${name}`)
      .then(response => response.json())
      .then(data => transformData(data))
      .then(data => {
        setData(data);
        // Receive parameters for filter
        if (menu['filter'].buttons.includes(method.filter) && method.filter !== 'all' && method.filterParam === 'all') {
          var newMenu = {...menu};
          newMenu[method.filter] = {
            buttons: ['all', ...data.map(item => item.label)],
            default: 'all',
          };
          setMenu(newMenu);
        }
      })
      .catch(error => console.error(error));
  };

  const [menu, setMenu] = useState(menuDefinition);

  return (
    <Container fluid>
      <Row>
        <Col>
          {getMenus(method).map(_menu => (
            <Menu
              key={_menu}
              onMethodChange={newMethod => handleMethodChange(newMethod, _menu)}
              buttons={menu[_menu]}
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
