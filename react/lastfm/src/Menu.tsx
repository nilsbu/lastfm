import React, { useState } from 'react';
import type { ButtonGroup } from './menus';
import 'bootstrap/dist/css/bootstrap.min.css';

interface MenuProps {
  onMethodChange: (newMethod: string) => void;
  buttons: ButtonGroup;
}

function Menu(props: MenuProps) {
  const [active, setActive] = useState(props.buttons.default);

  const handleClick = (name : string) => {
    setActive(name);
    props.onMethodChange(name);
  };

  return (
    <div className="bg-secondary p-3">
      <div className="row">
        <div className="col-auto">
          {props.buttons.buttons.map(button => (
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
