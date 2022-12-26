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
            <div className="container">
                <div className="row" style={{height: '10%'}}>
                    <Choices onSubmit={this.choose} type={Object.keys(OPTS)} page={this.state.page}/>
                </div>
                <Content page={this.state.page}/>
            </div>
        )
    }
}

function Content(props) {
    switch (props.page) {
    case "main":
        return (
            <div className="row" style={{height: '90%'}}>
                <div className="col-sm table-responsive" style={{height: '100%'}}><Charts name="year" param={YEAR}/></div>
                <div className="col-sm table-responsive" style={{height: '100%'}}><Charts name="fade" param="365"/></div>
                <div className="col-sm table-responsive" style={{height: '100%'}}><Charts name="fade" param="3653"/></div>
            </div>
        );
    default:
        return (
            <div className="row table-responsive" style={{height: '90%'}}>
                <ChosenCharts options={OPTS[props.page]} func={props.page} key={props.page} /> 
            </div>
        );
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
            <div>
                <div className="row" style={{height: '10%'}}>
                    <Choices onSubmit={this.choose} type={this.props.options} page={this.state.current}/>
                </div>
                <div className="col-sm table-responsive" style={{height: '100%'}} key={this.props.func+this.state.current}>
                    <Charts name={this.props.func} param={this.state.current}/>
                </div>
            </div>
    )}
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
        fetch(CMD[this.props.name](this.props.param))
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
            )
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
            <div className="charts-table">
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
            <nav aria-label="...">
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
