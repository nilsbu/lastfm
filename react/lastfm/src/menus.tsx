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
  'filter': {
    buttons: [
      { function: 'all', name: 'All' },
      { function: 'super', name: 'Super' },
      // { function: 'year', name: 'Year' },
    ],
    default: 'all',
  },
  'super': {
    buttons: [
      { function: 'all', name: 'All' },
      // will be filled in dynamically
    ],
    default: 'all',
  },
};

export type MenuChoice = {
  topLevel: string;
  functionParam: string; // there might be others in the future when we add intervals
  filter: string;
  filterParam: string;
};

export const getMenus = (method: MenuChoice): string[] => {
  const result: string[] = ['topLevel'];

  if (method.topLevel !== 'total') {
    result.push(method.topLevel);
  }
  
  result.push('filter');

  if (method.filter !== 'all') {
    result.push(method.filter);
  }

  return result;
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

export const transformMethod = (methodArray: MenuChoice): string[] => {
  var result: string[] = [methodArray.topLevel];

  if (methodArray.functionParam !== '') {
    result.push(methodArray.functionParam);
  }

  if (methodArray.filter !== 'all') {
    result.push(`by=${methodArray.filter}`);
    if (methodArray.filterParam !== 'all') {
      result.push(`name=${methodArray.filterParam}`);
    }
  }

  return result;
};
