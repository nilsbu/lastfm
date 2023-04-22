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

const normalizeMethod = (method: string[]): string[] => {
  // TODO: instead of this, the methods should be set correctly in the first place
  const result: string[] = method;
  if (method.length === 1) {
    // when only the top level is specified, add no filter
    result.push('all');
  }
  if (method.length === 3 && method[2] !== 'all') {
    // when only the top level and filter are specified, add no name
    result.push('all');
  }

  return result;
}

export const getMenus = (method: string[]): string[] => {
  method = normalizeMethod(method);
  const result: string[] = ['topLevel'];

  var i = 0;

  if (method.length > 0) {
    if (method[i] !== 'total') {
      result.push(method[i]);
      i++;
    }
    i++;

    result.push('filter');

    if (method.length > i && method[i] !== 'all') {
      result.push(method[i]);
      i++;
    }
    i++;
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

export const transformMethod = (methodArray: string[]) => {
  const index = methodArray.indexOf('super');
  if (index !== -1) {
    const [by, name, ...rest] = methodArray.slice(index);
    const preSuper = methodArray.slice(0, index);
    return [
      ...preSuper,
      `by=${by}`,
      name !== 'all' && name != null ? `name=${name}` : '',
      ...rest,
    ];
  } else {
    return methodArray.filter((element) => element !== 'all');
  }
};
