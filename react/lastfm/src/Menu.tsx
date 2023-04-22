import React, { useState } from 'react';
import type { ButtonGroup } from './menus';
import 'bootstrap/dist/css/bootstrap.min.css';

const BUTTONS_PER_PAGE = 5;

interface MenuProps {
  onMethodChange: (newMethod: string) => void;
  buttons: ButtonGroup;
}

function Menu(props: MenuProps) {
  const [active, setActive] = useState(props.buttons.default);
  const [startIndex, setStartIndex] = useState(0);

  const handleButtonClick = (name: string) => {
    setActive(name);
    props.onMethodChange(name);
  };

  const handlePrevClick = () => {
    const newStartIndex = Math.max(0, startIndex - BUTTONS_PER_PAGE);
    setStartIndex(newStartIndex);
  };

  const handleNextClick = () => {
    const newStartIndex = Math.min(
      props.buttons.buttons.length - BUTTONS_PER_PAGE,
      startIndex + BUTTONS_PER_PAGE
    );
    setStartIndex(newStartIndex);
  };

  const visibleButtons = props.buttons.buttons.slice(
    startIndex,
    startIndex + BUTTONS_PER_PAGE
  );

  return (
    <div className="bg-secondary p-3">
      <div className="row">
        <div className="col-auto">
          {startIndex > 0 && (
            <button className="btn btn-light mx-1" onClick={handlePrevClick}>
              &lt;
            </button>
          )}
          {visibleButtons.map((button) => (
            <button
              key={button.function}
              onClick={() => handleButtonClick(button.function)}
              className={`btn btn-light mx-1 ${
                active === button.function ? 'active' : ''
              }`}
            >
              {button.name}
            </button>
          ))}
          {startIndex + BUTTONS_PER_PAGE < props.buttons.buttons.length && (
            <button className="btn btn-light mx-1" onClick={handleNextClick}>
              &gt;
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

export default Menu;
