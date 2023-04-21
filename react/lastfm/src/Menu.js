import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

function Menu() {
  return (
    <div className="bg-secondary p-3 d-flex justify-content-between">
      <button className="btn btn-light">Button 1</button>
      <button className="btn btn-light">Button 2</button>
    </div>
  );
}

export default Menu;
