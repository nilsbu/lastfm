type Button = { function: string; name: string };

export type ButtonGroup = {
  buttons: Button[];
  default: string;
}

type MenuDefinition = {
  [key: string]: ButtonGroup;
};

export const buttons: MenuDefinition = {
  'topLevel': {buttons: [
    { function: 'total', name: 'Total' },
    { function: 'fade', name: 'Fade' },
  ], default: 'total'},
  'fade': {buttons:[
    { function: '30', name: '30' },
    { function: '365', name: '365' },
    { function: '1000', name: '1000' },
    { function: '3653', name: '3653' },
  ], default: '365'}
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
  