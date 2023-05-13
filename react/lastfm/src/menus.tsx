export type ButtonGroup = {
  buttons: string[];
  default: string;
}

type MenuDefinition = {
  [key: string]: ButtonGroup;
};

const currentYear = new Date().getFullYear();

const eraDelimiters: string[] = [
  '2007-05-05', // signed up to Last.fm
  '2012-02-17', // signed up to Simfy
  '2017-02-11', // started listening to Lady Gaga
  '2022-06-20', // moved to Spain
  '2030-12-31', 
];

export const menuDefinition: MenuDefinition = {
  'topLevel': {
    buttons: [
      'total',
      'fade',
      'period',
      'decades',
      'eras',
      'fademax',
    ],
    default: 'total'
  },
  'fade': {
    buttons:[
      '30',
      '365',
      '1000',
      '3653',
    ],
    default: '365'
  },
  'period': {
    buttons: Array.from({length: currentYear - 2006}, (_, i) => {
      const year = 2007 + i;
      return year.toString();
    }),
    default: currentYear.toString()
  },
  'decades': {
    buttons: [
      '2000s',
      '2010s',
      '2020s',
    ],
    default: '2020s'
  },
  'eras': {
    buttons: [
      'pre-streaming',
      'exploration',
      'pop',
      'international',
    ],
    default: 'international'
  },
  'fademax': {
    buttons:[
      '30',
      '365',
      '1000',
      '3653',
    ],
    default: '365'
  },
  'filter': {
    buttons: [
      'all',
      'super',
      'year',
      'country',
      'groups',
    ],
    default: 'all',
  },
  'super': {
    buttons: [
      'all',
      // will be filled in dynamically
    ],
    default: 'all',
  },
  'year': {
    buttons: [
      'all',
      // will be filled in dynamically
    ],
    default: 'all',
  },
  'country': {
    buttons: [
      'all',
      // will be filled in dynamically
    ],
    default: 'all',
  },
  'groups': {
    buttons: [
      'all',
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

  if (methodArray.topLevel === 'decades') {
    let decadeStart = methodArray.functionParam;
    let decadeEnd = (parseInt(decadeStart) + 10).toString() + 's';
  
    result = ['interval', `${decadeStart.slice(0, -1)}-01-01`, `${decadeEnd.slice(0, -1)}-01-01`];
  }
  else if (methodArray.topLevel === 'eras') {
    const eraIndex = menuDefinition.eras.buttons.indexOf(methodArray.functionParam);
    const eraStart = eraDelimiters[eraIndex];
    const eraEnd = eraDelimiters[eraIndex + 1];

    result = ['interval', eraStart, eraEnd];
    console.log(result);
  }
  else if (methodArray.functionParam !== '') {
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
