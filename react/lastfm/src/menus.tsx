type Button = { function: string; name: string };

export type ButtonGroup = {
  buttons: Button[];
  default: string;
}

type MenuDefinition = {
  [key: string]: ButtonGroup;
};

const currentYear = new Date().getFullYear();

export const menuDefinition: MenuDefinition = {
  'topLevel': {
    buttons: [
      { function: 'total', name: 'Total' },
      { function: 'fade', name: 'Fade' },
      { function: 'period', name: 'Period' },
      { function: 'super', name: 'Super'},
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
  },
  'super': {
    buttons: [
      { function: 'all', name: 'All' },
      // Will be filled in by the server
    ],
    default: 'all'
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
    case 'super':
      return ['topLevel', 'super'];
    default:
      return ['topLevel'];
  }
};

export const getQuery = (methodArray: string[]) => {
  let queryStringStarted = false;
  let result = '';
  for (let i = 0; i < methodArray.length; i++) {
    const element = methodArray[i];
    if (queryStringStarted || element.includes('=')) {
      if (!queryStringStarted) {
        result += '?';
        queryStringStarted = true;
      } else {
        result += '&';
      }
    } else if (i !== 0) {
      result += '/';
    }
    result += element;
  }
  return result;
};

export const transformMethod = (methodArray: string[]) => {
  if (methodArray[0] === 'super') {
    const [by, name, ...rest] = methodArray;
    return [
      'total',
      `by=${by}`,
      name !== 'all' ? `name=${name}` : '',
      ...rest,
    ].filter(Boolean);
  } else {
    return methodArray;
  }
};
