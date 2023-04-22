import React, { useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

function Menu(props) {
  const [active, setActive] = useState('total');

  const handleClick = (name) => {
    setActive(name);
    props.onMethodChange(name);
  };

  // button functionality
  const buttons = [
    { function: 'total', name: 'Total' },
    { function: 'fade/365', name: 'Fade 1y' },
    { function: 'fade/3653', name: 'Fade 10y' },
  ];

  return (
    <div className="bg-secondary p-3">
      <div className="row">
        <div className="col-auto">
          {buttons.map(button => (
            <button
              key={button.function}
              onClick={() => handleClick(button.function)}
              className={`btn btn-light mx-1 ${active === button.function ? 'active' : ''}`}
            >{button.name}</button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default Menu;
