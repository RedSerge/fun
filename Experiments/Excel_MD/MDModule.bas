'MultiDimensional Excel - a concept of a tensor-based table, where cells may contain multiple values in various dimensions.
'Written in VBA, compiled as an Add-In for Microsoft Excel 2019 (16.x).

'This part is provided by MS Excel VBA: module name
Attribute VB_Name = "MDModule"
'^ This part was provided by MS Excel VBA ^

'An instance of the custom handler class (described in 'SheetControl.cls')
Dim HandlerObject As New SheetControl

'Check if sheet with given name exists
Function IsSheet(ByVal Name As String)
    On Error Resume Next
    'If sheet exists, the function will return True 
    IsSheet = ActiveWorkbook.Sheets(Name).Name = Name
    'In case of error (no sheet with given name found), the string above is skipped
    'Due to instruction 'On Error Resume Next', the default (interpreted as 'False') value is returned as a result
End Function

'Form dimension name based on the number
Function MDName(ByVal Number As Integer)
    MDName = "MD(" + Trim(Str(Number)) + ")"
End Function

'Create (hidden or visible) page, if it does not exist yet
Sub CreatePage(ByVal Name As String, Optional ByVal Hidden As Boolean = False)
    If Not IsSheet(Name) Then ActiveWorkbook.Sheets.Add(Before:=ActiveWorkbook.Sheets(1)).Name = Name
    'Hide or select page based on the corresponding flag
    If Hidden Then
        ActiveWorkbook.Sheets(Name).Visible = xlSheetVeryHidden
    Else
        ActiveWorkbook.Sheets(Name).Select
    End If
End Sub

'Create dimension # [Number] (hidden page with corresponding name)
Sub CreateShadowPages(ByVal Number As Integer)
    For N = 1 To Number
        Name = MDName(N)
        CreatePage Name, True
    Next
End Sub

'Erase other dimensions after given number
Sub DestroyShadowPages(ByVal AfterNumber As Integer)
    'First dimension to erase 
    FromNumber = AfterNumber + 1
    'Prevent alert messages during sheet deletion process
    Application.DisplayAlerts = False
    Do While True
		'Name of page related to the dimension
        Name = MDName(FromNumber)
        'Delete page
        If IsSheet(Name) Then
            ActiveWorkbook.Sheets(Name).Visible = True
            ActiveWorkbook.Sheets(Name).Delete
        Else
			'If page does not exist, quit
            Exit Do
        End If
        'Otherwise, try to erase next dimension
        FromNumber = FromNumber + 1
    Loop
    'Restore default behaviour of alert messages
    Application.DisplayAlerts = True
End Sub

'Show dimension # [Number]
Sub SelectShadowPage(ByVal Number As Integer)
	'Check if hidden multidimensional storage of system info (inner memory) exists 
    If IsSheet("MD_Hidden") Then
		'Get current visible (selected) dimension from inner memory, correct the bounds
        InnerValue = PageNow()
        If InnerValue < 1 Then InnerValue = 1
        If Number < 1 Then Number = 1
        'Get name of the page related to the dimension selected according to inner memory  
        InnerName = MDName(InnerValue)
        'Get name of the page that is supposed to be selected
        Name = MDName(Number)
        'If both sheets exist, hide the former, show the latter;
        If IsSheet(Name) And IsSheet(InnerName) Then
            ActiveWorkbook.Sheets(InnerName).Visible = xlSheetVeryHidden
            ActiveWorkbook.Sheets(Name).Visible = True
            'also, select the new selected one and update inner memory
            ActiveWorkbook.Sheets(Name).Select
            SetPageNow (Number)
        End If
    End If
End Sub

'Find last multidimensional page (proceed with search till page does not exist)
'0 = no multidimensional pages found
Function LastPage()
    LastPage = 0
    While IsSheet(MDName(LastPage + 1))
        LastPage = LastPage + 1
    Wend
End Function

'Get current visible (selected) dimension (read from inner memory)
Function PageNow()
    If IsSheet("MD_Hidden") Then PageNow = ActiveWorkbook.Sheets("MD_Hidden").Cells(1, 1).Value
End Function

'Set current visible (selected) dimension (write to inner memory)
Sub SetPageNow(ByVal Count As Integer)
    If IsSheet("MD_Hidden") Then ActiveWorkbook.Sheets("MD_Hidden").Cells(1, 1).Value = Count
End Sub

'Perform dimension increment or decrement based on the direction 
Sub ChangePage(ByVal Direction As Integer)
    If IsSheet("MD_Hidden") Then
        NewPage = PageNow() + Direction
        'If supposed page is out of bounds, make it equal to bound opposite to the broken one
        If NewPage < 1 Then NewPage = LastPage
        If Not IsSheet(MDName(NewPage)) Then NewPage = 1
        'Select that page
        SelectShadowPage (NewPage)
    End If
End Sub

'Action procedure for button 'Next dimension'
Sub NextPage()
    ChangePage (1)
End Sub

'Action procedure for button 'Previous dimension'
Sub PrevPage()
    ChangePage (-1)
End Sub

'Action procedure for button 'Configuration';
'this process also automatically restructures current workbook as multidimensional
Sub SetConfig()
	'Create inner memory if it doesn't exist (form the multidimensional structure in its basic form)
    If Not IsSheet("MD_Hidden") Then
        CreatePage "MD_Hidden", True
        DestroyShadowPages (1)
        CreateShadowPages (1)
        SetPageNow (1)
    End If
    'Create the Configuration Page if it doesn't exist yet
    CreatePage "MD_Config", True
    'Get current visible page and overall count of the multidimensional pages (hidden sheets)
    Current = PageNow()
    Count = LastPage()
    'Show, select and clear the Configuration Page 
    ActiveWorkbook.Sheets("MD_Config").Visible = True
    Sheets("MD_Config").Select
    Cells.Clear
    Cells.MergeCells = False
    'Reset the width of the 1st column
    Columns("A:A").ColumnWidth = 25
    'Fill the page with appropriate values
    Cells(1, 1).Value = "Visible dimension #:"
    Cells(1, 2).Value = Current
    Cells(2, 1).Value = "Dimension count:"
    Cells(2, 2).Value = Count
    'Prepare the sorting zone
    Range("E1:F1").MergeCells = True
    With Range("E1:F2")
        .HorizontalAlignment = xlCenter
        .VerticalAlignment = xlCenter
    End With
    Cells(1, 5).Value = "Sort"
    Cells(2, 5).Value = "Old"
    Cells(2, 6).Value = "New"
    'Start propagating the formula for each dimension number in the left sorting column ("Old")
    Cells(3, 5).Value = 1
    If Count > 1 Then Cells(4, 5).FormulaR1C1 = "=R[-1]C+1"
    If Count > 2 Then Cells(4, 5).AutoFill Destination:=Range("E4:E" + Trim(Str(Count + 2))), Type:=xlFillDefault
    'Create the 'LAUNCH' cell ("button")
    Range("A6:A7").MergeCells = True
    With Range("A6:A7")
        .HorizontalAlignment = xlCenter
        .VerticalAlignment = xlCenter
        .Interior.ThemeColor = xlThemeColorAccent1
        .Font.ThemeColor = xlThemeColorDark2
        .Font.Bold = True
        .Font.Size = 14
        .Value = "LAUNCH"
    End With
    'Select default cell and reset handler
    Cells(1, 2).Select
    Set HandlerObject.sheetevent = Nothing
    Set HandlerObject.sheetevent = ActiveWorkbook.Sheets("MD_Config")    
End Sub

'Action procedure for button 'Store selected area'
Sub Selection_Store()
	'Save address of the selected range to inner memory
    If IsSheet("MD_Hidden") Then ActiveWorkbook.Sheets("MD_Hidden").Cells(2, 1).Value = Selection.Address
End Sub

'Action procedure for button 'No selection'
Sub Selection_Empty()
	'Clean inner memory related to saved range address
    If IsSheet("MD_Hidden") Then ActiveWorkbook.Sheets("MD_Hidden").Cells(2, 1).Value = ""
End Sub

'Action procedure for button 'Show stored area'
Sub Selection_Recover()
    If IsSheet("MD_Hidden") Then
		'Select saved range address; if nothing is saved, show the corresponding message
        selection_stored = ActiveWorkbook.Sheets("MD_Hidden").Cells(2, 1).Value
        If selection_stored = "" Then MsgBox ("(Empty)") Else ActiveSheet.Range(selection_stored).Select
    End If
End Sub

'Multidimensional copying and/or deletion of the cells
'(operation selection is based on the provided function arguments) 
Sub Selection_Transform(ByVal paste As Boolean, ByVal clean As Boolean)
    If IsSheet("MD_Hidden") Then
		'Store the address of current selection range, in case of copying operation
        If paste Then selection_passed = Selection.Cells(1, 1).Address
        'Load the saved selection range from inner memory
        selection_stored = ActiveWorkbook.Sheets("MD_Hidden").Cells(2, 1).Value
        '...If it is actually saved, that's it
        If selection_stored <> "" Then
			'Go through each dimensional sheet and perform required operations
            Index = 1
            While IsSheet(MDName(Index))
                With ActiveWorkbook.Sheets(MDName(Index))
                'Copy...
                If paste Then .Range(selection_stored).Copy .Range(selection_passed)
                '...and/or Delete info related to the corresponding areas
                If clean Then .Range(selection_stored).Clear
                End With
            Index = Index + 1
            Wend
        End If
    End If
End Sub

'Action procedure for button 'Erase stored area'
Sub Action_Erase()
    Selection_Transform False, True
End Sub

'Action procedure for button 'Fill with stored area'
Sub Action_Copy()
    Selection_Transform True, False
End Sub

'Action procedure for button 'Move from stored area'
Sub Action_Cut()
    Selection_Transform True, True
End Sub