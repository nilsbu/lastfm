const YEAR = new Date().getFullYear();

const PERIODS = [
    ["2007-05-05", "2012-02-17"],
    ["2012-02-17", "2017-02-11"],
    ["2017-02-11", "2022-06-20"],
    ["2022-06-20", "2099-12-31"],
];

const CMD = {
    "year":     (v) => `/json/print/period/${v}`,
    "fromYear": (v) => `/json/print/total?by=year&name=${v}`,
    "fade":     (v) => `/json/print/fade/${v}`,
    "super":    (v) => `/json/print/total?by=super&name=${v}`,
    "period":   (v) => `/json/print/interval/${PERIODS[v][0]}/${PERIODS[v][1]}`,
};

const YEARS = Array.from({length: YEAR - 2007 + 1}, (x, i) => i + 2007); // TODO fix init

const SUPERS = ["rock", "metal", "pop", "electronic", "hip-hop", "folk", "reggae", "classical", "jazz"];

const OPTS = {
    "main": [],
    "buffet": [],
    "fromYear": YEARS,
    "year": YEARS,
    "super": SUPERS,
    "fade": [30, 365, 1000, 3653],
    "period": [0, 1, 2, 3],
};

class Dashboard extends React.Component {
    constructor(props) {
        super(props);
        this.state = { 
            page: Object.keys(OPTS)[0],
        };
    }

    choose = (page) => {
        this.setState(Object.assign({}, this.state, {
            page: page,
        }));
    }

    render() {
        return (
            <div className="container main">
                <div className="row row-cols-2" >
                    <Choices onSubmit={this.choose} type={Object.keys(OPTS)} page={this.state.page}/>
                </div>
                <Content page={this.state.page}/>
            </div>
        );
    }
}

function Content(props) {
    switch (props.page) {
    case "main":
        return (
            <div className="row row-body" style={{display: "block"}}>
                <Charts func={CMD["year"]} param={YEAR}/>
                <Charts func={CMD["fade"]} param="365"/>
                <Charts func={CMD["fade"]} param="3653"/>
            </div>
        );
    case "buffet":
        return (
            <div className="row row-body" style={{display: "block"}}>
                <Buffet />
            </div>
        );
    default:
        return (
            // <div className="row row-body table-responsive">
                <ChosenCharts options={OPTS[props.page]} func={CMD[props.page]} key={props.page} /> 
            // </div>
        );
    }
}

class Buffet extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            base: "total",
            params: {},
            jsxParams: null,
        };

        this.chooseBase = (page) => {
            this.setState(Object.assign({}, this.state, {
                base: page,
                params: {},
            }));
        };

        this.setParams = (params) => this.setState(Object.assign({}, this.state, {
            params: params,
        })); 
    }

    getFunc() {
        // TODO: is the closure here needed?
        var str = "/json/print/" + this.state.base;
        var params = {...this.state.params};

        switch (this.state.base) {
        case "fade":
        case "period":
            str += `/${params.p0}`;
            break;
        case "interval":
            str += `/${params.p0}/${params.p1}`;
            break;
        }
        delete params.p0;
        delete params.p1;
        
        var first = true;
        for (const [key, value] of Object.entries(params)) {
            if (first) {
                str += "?";
                first = false;
            } else {
                str += "&";
            }
            str += `${key}=${value}`;
        }

        return function (param) {return str;}
    }

    render() {
        return (
            <div key={`buffet-${this.getFunc()()}`}>
                <div className="row">
                    <Choices onSubmit={this.chooseBase} type={["total", "fade", "period", "interval"]} page={this.state.base}/>
                </div>
                <Params base={this.state.base} cb={this.setParams} />
                <Charts func={this.getFunc()} param="" />
            </div>
        );
    }
}

class Params extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            params: {},
            lastBase: props.base,
        };

        this.set = (name, value) => {
            this.setState(Object.assign({}, this.state, {
                params: Object.assign({}, this.state.params, {
                    [name]: value,
                }),
            }));
        }; 
    }

    componentDidUpdate(prevProps) {
        if (prevProps.base !== this.props.base) {
            this.setState(Object.assign({}, this.state, {
                params: {},
            }));
            this.props.cb(this.state.params);
        }
    }

    render() {
        var titles = [];

        switch (this.props.base) {
        case "total":
            break;
        case "fade":
            titles = ["half-life"];
            break;
        case "period":
            titles = ["period"];
            break;
        case "interval":
            titles = ["begin", "end"];
            break;
        }

        if (titles.length == 0) {
            return (<div/>);
        } else {
            return (
                <div className="input-group bg-dark">
                    {
                        titles.map((opt, i) => [
                            (<span className="input-group-text bg-dark" key={"p-span-" + opt} >{opt}</span>), 
                            (<input type="text" key={"p-input-" + opt} className="form-control bg-dark" placeholder="" aria-label={opt} aria-describedby="basic-addon1" onChange={(v) => this.set(`p${i}`, v.target.value)} />)]
                        )
                    }
                    <button type="button" className="btn btn-outline-secondary bg-dark" onClick={() => this.props.cb(this.state.params)} >Confirm</button>
                </div>
            );
        }
    }
}

class ChosenCharts extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            current: props.options[0],
        };

        this.choose = (page) => this.setState(Object.assign({}, this.state, {
            current: page,
        }));
    }

    render() {
        return (
            <div className="row-body row container" style={{display: "block"}} key={`cc-${this.props.options}-${this.state.current}`}>
                <Choices onSubmit={this.choose} type={this.props.options} page={this.state.current}/>
                <Charts func={this.props.func} param={this.state.current}/>
            </div>
        );
    }
}

class Charts extends React.Component {
    constructor(props) {
        super(props);
        this.state = { 
            data: null,
            isLoaded: false,
        };
    }

    componentDidMount() {
        fetch(this.props.func(this.props.param))
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState(Object.assign({}, this.state, {
                        isLoaded: true,
                        data: this.reshapeData(result),
                    }));
                },
                (error) => {
                    this.setState({
                    isLoaded: true,
                    error
                    });
                }
            );
    }

    reshapeData(raw) {
        var result = {"elems": []};
        var data = raw.data.datasets[0].data;
        for (var i = 0; i < data.length; i++) {
            result.elems.push({"name": raw.data.labels[i], "value": data[i]});
        }
        return result;
    }

    render() {
        if (!this.state.data) {
            return <div className="charts-table"/>
        }
        return (
            <div className="col-sm table-responsive charts-table row">
                <table className="table table-striped bg-dark text-white"><tbody>
                    {this.state.data.elems.map((elem, i) => <Line key={elem.name} idx={i} elem={elem}/>)}
                </tbody></table>
            </div>
        );
    }
}

class Choices extends React.Component {
    constructor(props) {
        super(props);
        this.state = { 
            options: props.type,
            page: props.page,
            onSubmit: props.onSubmit,
        };
    }

    render() {
        return (
            <nav aria-label="..." className="row">
                <ul className="pagination bg-dark">
                    {
                        this.state.options.map((opt) => (
                            <li className="page-item bg-dark" key={opt}>
                                <div className="page-link bg-dark" onClick={() => this.state.onSubmit(opt)}>{opt}</div>
                            </li>
                        ))
                    }
                </ul>
            </nav>
        );
    }
}

function Line(props) {
    return (
        <tr>
            <td>{props.idx+1}</td>
            <td>{props.elem.name}</td>
            <td>{props.elem.value.toFixed(2)}</td>
        </tr>
    );
}

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<Dashboard />);
