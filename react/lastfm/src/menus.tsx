export type Button = { function: string; name: string };

type ButtonGroup = {
  [key: string]: Button[];
};

export const buttons: ButtonGroup = {
    'topLevel': [
      { function: 'total', name: 'Total' },
      { function: 'fade', name: 'Fade' },
    ],
    'fade': [
      { function: '365', name: '365' },
      { function: '3653', name: '3653' },
    ]
  };
  
  export const getMenus = (topLevelFunction : string) => {
    switch (topLevelFunction) {
      case 'total':
        return ['topLevel'];
      case 'fade':
        return ['topLevel', 'fade'];
      default:
        return ['topLevel'];
    }
  };
  