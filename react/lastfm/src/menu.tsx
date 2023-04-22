export const buttons = {
    'topLevel': [
      { function: 'total', name: 'Total' },
      { function: 'fade', name: 'Fade' },
    ],
    'fade': [
      { function: '365', name: '365' },
      { function: '3653', name: '3653' },
    ]
  };
  
  export const getMenus = (topLevelFunction) => {
    switch (topLevelFunction) {
      case 'total':
        return ['topLevel'];
      case 'fade':
        return ['topLevel', 'fade'];
      default:
        return ['topLevel'];
    }
  };
  