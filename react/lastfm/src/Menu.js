import React, { useState } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

function Menu(props) {
  const [active, setActive] = useState('total');

  const handleClick = (name) => {
    setActive(name);
    props.onMethodChange(name);
  };

  return (
    <div className="bg-secondary p-3 d-flex justify-content-between">
      <button onClick={() => handleClick('total')} className={active === 'total' ? 'active' : ''}>Total</button>
      <button onClick={() => handleClick('fade')} className={`btn btn-light ${active === 'fade' ? 'active' : ''}`}>Fade</button>
    </div>
  );
}

export default Menu;
