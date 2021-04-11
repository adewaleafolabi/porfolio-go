

export const interpolateColors = (dataLength, colorScale, colorRangeInfo) => {
    const {colorStart, colorEnd} = colorRangeInfo;
    const colorRange = colorEnd - colorStart;
    const intervalSize = colorRange / dataLength;
    let i, colorPoint;
    const colorArray = [];

    for (i = 0; i < dataLength; i++) {
        colorPoint = calculateColorPoint(i, intervalSize, colorRangeInfo);
        colorArray.push(colorScale(colorPoint));
    }

    return colorArray;
}


const  calculateColorPoint = (i, intervalSize, colorRangeInfo)=> {
    const {colorStart, colorEnd, useEndAsStart} = colorRangeInfo;
    return (useEndAsStart
        ? (colorEnd - (i * intervalSize))
        : (colorStart + (i * intervalSize)));
}