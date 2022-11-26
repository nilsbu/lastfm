const CMD = {
    "year": `/json/print/period/${new Date().getFullYear()}`,
    "f365": "/json/print/fade/365",
    "f3653": "/json/print/fade/3653",
}

function Dashboard(props) {
    return (
        
        <div className="container"><div className="row">
            <div className="col"><Charts name="year"/></div>
            <div className="col"><Charts name="f365"/></div>
            <div className="col"><Charts name="f3653"/></div>
        </div></div>
    )
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
        fetch(CMD[this.props.name])
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
                <table><tbody>
                    {this.state.data.elems.map((elem) => <Line key={elem.name} elem={elem}/>)}
                </tbody></table>
            </div>
        );
  }
}

function Line(props) {
    return (
        <tr>
            <th>{props.elem.name}</th>
            <td>{props.elem.value.toFixed(2)}</td>
        </tr>
    );
}

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<Dashboard />);
