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
    { function: 'period', name: 'Period' },
  ], default: 'total'},
  'fade': {buttons:[
    { function: '30', name: '30' },
    { function: '365', name: '365' },
    { function: '1000', name: '1000' },
    { function: '3653', name: '3653' },
  ], default: '365'},
  'period': {buttons:[
    { function: '2007', name: '2007' },
    { function: '2008', name: '2008' },
    { function: '2009', name: '2009' },
    { function: '2010', name: '2010' },
    { function: '2011', name: '2011' },
    { function: '2012', name: '2012' },
    { function: '2013', name: '2013' },
    { function: '2014', name: '2014' },
    { function: '2015', name: '2015' },
    { function: '2016', name: '2016' },
    { function: '2017', name: '2017' },
    { function: '2018', name: '2018' },
    { function: '2019', name: '2019' },
    { function: '2020', name: '2020' },
    { function: '2021', name: '2021' },
    { function: '2022', name: '2022' },
    { function: '2023', name: '2023' },
  ], default: '2023'}
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
  