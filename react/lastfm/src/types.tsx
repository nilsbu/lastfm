// Here is where we define the types for our application

interface JSONElement {
    title: string;
    value: number;
    prevPos?: number;
    prevValue?: number;
}

// type that we get as JSON. There is more because it's also used for the chart.
interface JSONData {
    chart: {
        data: JSONElement[];
    };
    precision: number;
};
