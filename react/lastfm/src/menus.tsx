type Button = { function: string; name: string };

export type ButtonGroup = {
  buttons: Button[];
  default: string;
}

type MenuDefinition = {
  [key: string]: ButtonGroup;
};

const currentYear = new Date().getFullYear();

export const buttons: MenuDefinition = {
  'topLevel': {
    buttons: [
      { function: 'total', name: 'Total' },
      { function: 'fade', name: 'Fade' },
      { function: 'period', name: 'Period' },
    ],
    default: 'total'
  },
  'fade': {
    buttons:[
      { function: '30', name: '30' },
      { function: '365', name: '365' },
      { function: '1000', name: '1000' },
      { function: '3653', name: '3653' },
    ],
    default: '365'
  },
  'period': {
    buttons: Array.from({length: currentYear - 2006}, (_, i) => {
      const year = 2007 + i;
      return { function: year.toString(), name: year.toString() };
    }),
    default: currentYear.toString()
  }
};
  
  export const getMenus = (topLevelFunction : string) => {
    switch (topLevelFunction) {
      case 'total':
        return ['topLevel'];
      case 'fade':
        return ['topLevel', 'fade'];
      case 'period':
        return ['topLevel', 'period'];
      default:
        return ['topLevel'];
    }
  };
  