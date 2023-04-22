import React, { useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

function Menu(props) {
  const [active, setActive] = useState('total');

  const handleClick = (name) => {
    setActive(name);
    props.onMethodChange(name);
  };

  return (
    <div className="bg-secondary p-3">
      <div className="row">
        <div className="col-auto">
          {props.buttons.map(button => (
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
